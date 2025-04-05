sudo cp $HOME/jobby/jobby.service /etc/systemd/system/jobby.service
sudo systemctl enable jobby.service
sudo systemctl start jobby.service
sudo systemctl restart jobby.service
