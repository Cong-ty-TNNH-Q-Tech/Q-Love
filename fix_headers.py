import glob
import os

files = glob.glob('backend/server/**/*.go', recursive=True)
for f in files:
    with open(f, 'r', encoding='utf-8') as file:
        content = file.read()
    
    if '// Copyright (c) 2026 Q-Tech. All rights reserved.\n// Licensed under the GNU AGPLv3 License.\n' in content:
        content = content.replace(
            '// Copyright (c) 2026 Q-Tech. All rights reserved.\n// Licensed under the GNU AGPLv3 License.\n',
            '// Copyright 2026 Q-Tech Team\n// Licensed under the GNU AGPLv3 License.\n// See LICENSE file in the project root for full license information.\n'
        )
        with open(f, 'w', encoding='utf-8') as file:
            file.write(content)
        print(f"Fixed {f}")
