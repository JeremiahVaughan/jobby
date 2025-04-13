set -e
GOOS=linux GOARCH=arm64 go build -o ./deploy/app
rsync -avz --delete -e ssh ./deploy/ ./ui deploy.target:$HOME/jobby
ssh deploy.target "$HOME/jobby/remote-deploy.sh"

