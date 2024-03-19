import fnmatch
import os
import sys
import time
import json

def load_gitignore_patterns(gitignore_path):
    """Load and return .gitignore patterns."""
    patterns = ['.git/']  # Explicitly ignore .git directory
    try:
        with open(gitignore_path, 'r') as f:
            for line in f:
                line = line.strip()
                if line and not line.startswith('#'):
                    patterns.append(line)
    except FileNotFoundError:
        print(f"No .gitignore file found at {gitignore_path}. Continuing without it.")
    return patterns

def should_ignore(file, patterns, dir_path):
    """Determine if a file should be ignored based on .gitignore patterns and explicit ignores."""
    rel_path = os.path.relpath(file, dir_path)
    for pattern in patterns:
        if fnmatch.fnmatch(rel_path, pattern) or fnmatch.fnmatch(file, pattern):
            return True
    return False

def generate_file_data(dir_path, gitignore_patterns):
    """Generate file data not matching .gitignore patterns, handling UnicodeDecodeErrors."""
    files_data = []
    for root, dirs, files in os.walk(dir_path, topdown=True):
        # Filter out ignored directories
        dirs[:] = [d for d in dirs if not should_ignore(os.path.join(root, d), gitignore_patterns, dir_path)]
        for name in files:
            if not should_ignore(os.path.join(root, name), gitignore_patterns, dir_path):
                file_path = os.path.join(root, name)
                rel_path = os.path.relpath(file_path, dir_path)
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        content = f.read()
                        files_data.append({
                            "path": rel_path,
                            "content": content,
                            "filename": name
                        })
                except UnicodeDecodeError:
                    print(f"Skipping binary or non-UTF-8 file: {file_path}")
    return files_data

def write_json_file(output_file_path, files_data):
    """Write the files data as a JSON array to a file."""
    with open(output_file_path, 'w', encoding='utf-8') as output_file:
        json.dump(files_data, output_file, ensure_ascii=False, indent=4)

def main():
    if len(sys.argv) < 2:
        print("Usage: python script.py <directory>")
        sys.exit(1)
    
    dir_path = sys.argv[1]  # Take the directory as a command-line argument
    current_timestamp = int(time.time())
    output_file_name = f"concatenated_{current_timestamp}.json"
    output_file_path = os.path.join(dir_path, output_file_name)
    gitignore_path = os.path.join(dir_path, '.gitignore')
    
    gitignore_patterns = load_gitignore_patterns(gitignore_path)
    files_data = generate_file_data(dir_path, gitignore_patterns)
    write_json_file(output_file_path, files_data)
    print(f"Files have been saved into {output_file_path}")

if __name__ == "__main__":
    main()
