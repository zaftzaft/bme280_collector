service=bme280-collector

systemctl is-enabled $service
if [ $? -eq 0 ]; then
  systemctl stop $service
fi


go build

install -Dm755 bme280_collector /usr/bin/bme280_collector

if [ -d /usr/lib/systemd/system/ ]; then
  unit_dir=/usr/lib/systemd/system
else
  unit_dir=/etc/systemd/system
fi

install -Dm644 systemd/$service.service $unit_dir/$service.service
install -Dm644 systemd/$service.timer $unit_dir/$service.timer


systemctl daemon-reload
systemctl enable $service.timer
systemctl start $service.timer

