grep -rnw . --exclude 'moment.js' --exclude 'todo.sh' --exclude-dir 'vendor' --exclude-dir 'tmp' --exclude-dir '.git' --exclude-dir '.hg' --exclude-dir 'dist' -e 'TODO'
