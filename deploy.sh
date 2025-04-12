set -e
GOOS=linux GOARCH=arm64 go build -o ./deploy/app
rsync -avz --delete -e ssh ./deploy/ ./ui pi3:$HOME/jobby
ssh pi3 "$HOME/jobby/remote-deploy.sh"

