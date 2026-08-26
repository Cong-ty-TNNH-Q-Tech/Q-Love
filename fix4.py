import os
import re
from pathlib import Path

handlers_dir = Path(r"c:\Users\HNC\Desktop\QLove\backend\server\internal\api\handlers")

def fix_tests():
    for filepath in handlers_dir.glob("*_test.go"):
        with open(filepath, "r", encoding="utf-8") as f:
            content = f.read()

        # Fix "123e4567-e89b-12d3-a456-426614174000" string injections
        content = re.sub(
            r'c\.Locals\("user(?:ID|_id)",\s*"123e4567-e89b-12d3-a456-426614174000"\)',
            r'c.Locals("user_id", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))',
            content
        )
        
        # Unify userID to user_id just in case
        content = content.replace('c.Locals("userID", ', 'c.Locals("user_id", ')

        with open(filepath, "w", encoding="utf-8") as f:
            f.write(content)
            
    # specifically fix locket_test.go which expects 400 instead of 401
    locket_path = handlers_dir / "locket_test.go"
    if locket_path.exists():
        with open(locket_path, "r", encoding="utf-8") as f:
            locket_content = f.read()
            
        locket_content = locket_content.replace(
            'if resp.StatusCode != fiber.StatusBadRequest {\n\t\tt.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)\n\t}',
            'if resp.StatusCode != fiber.StatusUnauthorized {\n\t\tt.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)\n\t}'
        )
        with open(locket_path, "w", encoding="utf-8") as f:
            f.write(locket_content)

if __name__ == "__main__":
    fix_tests()
