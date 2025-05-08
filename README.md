# wrapl

**wrapl** is a simple REPL wrapper for any shell command with a `{{}}` placeholder. It lets you interactively run variations of a shell command by replacing `{{}}` with your input at runtime.

## Features

- Live REPL with history and line editing (via [chzyer/readline](https://github.com/chzyer/readline))
- Substitutes `{{}}` in the base command template with your input
- Great for working with databases, CLI tools, or any repetitive shell command

## Installation

```bash
go install github.com/KishorPokharel/wrapl@latest
```

## Usage

```bash
wrapl "npx wrangler d1 execute my-database --remote --env uat --command \"{{}}\""
```
This launches an interactive shell. Type SQL queries, and wrapl will run the full command with your input substituted into {{}}.

Example:

```bash
$ wrapl "npx wrangler d1 execute my-database --remote --env uat --command \"{{}}\""

REPL started. Type 'exit' to quit.
> SELECT * FROM notes;
Result:
+----+----------+
| id | content  |
+----+----------+
| 1  | hello    |
| 2  | world    |
+----+----------+
> INSERT INTO notes (content) VALUES ('wrapl is cool');
Success
> exit
```
