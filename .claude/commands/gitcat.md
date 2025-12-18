---
description: Load a git repository into context using gitcat
argument-hint: <repo-url> [--md|--jsonl|--text] [--to-file <filename>]
---

Load a git repository into Claude's context using the gitcat tool.

**Arguments from user:** $ARGUMENTS

**Supported parameters:**
- `<repo-url>`: Repository URL or path (required) - HTTPS, SSH, or local path
- `--md`: Use markdown format (default if no format specified)
- `--jsonl`: Use JSONL format
- `--text`: Use text format
- `--to-file <filename>`: Save to file instead of loading into context

**Steps:**

1. **Parse arguments** from $ARGUMENTS:
   - Extract repo URL (first non-flag argument)
   - Detect format flag: `--md`, `--jsonl`, or `--text` (default: md)
   - Extract filename if `--to-file` is present
   - If no repo URL provided, ask the user for it

2. **Determine format:**
   - If `--md` present: format = "md"
   - If `--jsonl` present: format = "jsonl"
   - If `--text` present: format = "text"
   - Otherwise: format = "md" (default)

3. **Determine output destination:**
   - If `--to-file <filename>` present: use `-out <filename>` flag
   - Otherwise: write to `/tmp/gitcat-context` and read it back

4. **Build and execute gitcat command:**
   ```bash
   # If saving to user file:
   gitcat -fmt <format> -tmp -out <filename> <repo-url>

   # If loading into context:
   gitcat -fmt <format> -tmp -out /tmp/gitcat-context <repo-url>
   ```

   Use `-tmp` flag for remote repos to auto-cleanup the clone.

5. **Load the output:**
   - If saved to file: inform user where the file was saved (e.g., "myrepo.md")
   - If loading to context: read `/tmp/gitcat-context.<format>` using the Read tool

6. **Analyze and summarize:**
   - Count total files loaded
   - Identify file types/extensions
   - List main directories/packages
   - Identify programming languages
   - Note special files (README, LICENSE, config)

7. **Respond to user:**
   - Provide structured summary:
     - Repository information
     - File statistics
     - Technology stack detected
     - Project structure overview
   - If loaded into context, say: "Repository loaded! I can now help you understand, analyze, or modify this codebase. What would you like to know?"

**Error handling:**
- If gitcat fails, explain the error
- Check if gitcat is in PATH
- Verify URL format for remote repos
- Remind about SSH keys for SSH repos

**Example usages:**
```
/gitcat https://github.com/user/repo.git
/gitcat https://github.com/user/repo.git --md
/gitcat https://github.com/user/repo.git --jsonl
/gitcat git@github.com:user/repo.git --md --to-file myrepo
/gitcat . --text
/gitcat /path/to/local/repo --to-file project-context
```
