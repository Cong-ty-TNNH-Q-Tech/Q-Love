import os
import re
from pathlib import Path

handlers_dir = Path(r"c:\Users\HNC\Desktop\QLove\backend\server\internal\api\handlers")

# Regex to find c.Locals("userID", uuid.New().String())
# or c.Locals("user_id", uuid.New().String())
# or c.Locals("user_id", uuid.NewString())
pattern1 = re.compile(r'c\.Locals\("user(?:ID|_id)",\s*uuid\.New\(\)\.String\(\)\)')
pattern2 = re.compile(r'c\.Locals\("user(?:ID|_id)",\s*uuid\.NewString\(\)\)')

# Wait, what about other random strings? Let's also unify the key to "user_id".
def fix_tests():
    for filepath in handlers_dir.glob("*_test.go"):
        with open(filepath, "r", encoding="utf-8") as f:
            content = f.read()

        new_content = pattern1.sub('c.Locals("user_id", uuid.New())', content)
        new_content = pattern2.sub('c.Locals("user_id", uuid.New())', new_content)
        
        # also there might be c.Locals("userID", uuid.New()) without string(), just unify the key
        new_content = new_content.replace('c.Locals("userID", uuid.New())', 'c.Locals("user_id", uuid.New())')
        
        if new_content != content:
            with open(filepath, "w", encoding="utf-8") as f:
                f.write(new_content)
            print(f"Fixed {filepath.name}")

if __name__ == "__main__":
    fix_tests()
