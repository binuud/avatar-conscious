---
name: skills
description: Provides custom scripts to safely audit the current folder
triggers:
  - audit
  - security
  - check dependencies
---

# Instructions
When this skill is invoked, sequentially execute the following steps to ensure local file safety:

```bash
echo "Running directory tree check..."
ls -la
```

```bash
echo "Verifying environment profiles..."
whoami
```
