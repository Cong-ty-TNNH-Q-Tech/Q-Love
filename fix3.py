import os
import re

files_to_fix = [
    r"c:\Users\HNC\Desktop\QLove\backend\server\internal\api\handlers\court_handler.go",
    r"c:\Users\HNC\Desktop\QLove\backend\server\internal\api\handlers\dating_contract_handler.go",
    r"c:\Users\HNC\Desktop\QLove\backend\server\internal\api\handlers\ex_rating_handler.go",
]

for filepath in files_to_fix:
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()

    # Pattern: 
    # varStr := c.Locals("userID").(string) (or user_id)
    # varID, err := uuid.Parse(varStr)
    # if err != nil { return ... }
    pattern = re.compile(
        r'(\w+?)(?:Str|Val)?(?:\s*,\s*ok)?\s*:=\s*c\.Locals\("user(?:ID|_id)"\)\.\(string\)\s*\n'
        r'(?:\s*if\s*!ok\s*(?:\|\|\s*\1(?:Str|Val)?\s*==\s*"")?\s*\{\s*\n(?:.|\n)*?\s*\}\s*\n)?'
        r'\s*(\w+),\s*err\s*:=\s*uuid\.Parse\(\1(?:Str|Val)?\)\s*\n'
        r'\s*if\s*err\s*!=\s*nil\s*\{\s*\n'
        r'(?:.|\n)*?'
        r'\s*\}\s*\n'
    )
    def repl(m):
        uuid_var = m.group(2)
        return f'{uuid_var}, ok := c.Locals("user_id").(uuid.UUID)\n\tif !ok {{\n\t\treturn c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{{"error": "Unauthorized"}})\n\t}}\n'
    content = pattern.sub(repl, content)

    with open(filepath, "w", encoding="utf-8") as f:
        f.write(content)

print("Done fixing 3 files")
