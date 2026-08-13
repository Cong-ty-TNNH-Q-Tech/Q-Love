import re
import sys

def fix_api():
    with open('docs/api.yaml', 'r', encoding='utf-8') as f:
        content = f.read()

    # Find the conflict block
    # It looks like:
    # <<<<<<< HEAD
    # (locket stuff)
    # =======
    # (clans stuff)
    # >>>>>>> hash
    
    # We want to replace it with:
    # (locket stuff)
    # (clans stuff)

    pattern = re.compile(r'<<<<<<< HEAD\n(.*?)=======\n(.*?)>>>>>>> [a-f0-9]+', re.DOTALL)
    
    match = pattern.search(content)
    if match:
        new_content = content[:match.start()] + match.group(1) + match.group(2) + content[match.end():]
        with open('docs/api.yaml', 'w', encoding='utf-8') as f:
            f.write(new_content)
        print("Fixed conflict in api.yaml")
    else:
        print("No conflict marker found")

fix_api()
