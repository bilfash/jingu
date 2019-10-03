# Jingu


### Description
Download attachment from email

### Usage
```bash
 jingu path_to_config.yml
```

### Example config

```yaml
host: imap.google.com
port: 993
username: jingu@gmail.com
password: jingu
mailbox: INBOX
subjects: subject1,subject2,subject3
file_pattern: .txt$
sink_folder: .
```

### Installation

```bash
go install github.com/bilfash/jingu
```