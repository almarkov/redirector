# redirector

redirects user requests to different hosts using ip information  
geoip info from https://ip-api.com/  

REDIRECTOR
```
cd /var/projects/
mkdir redirector
git clone https://github.com/almarkov/redirector.git
go build redirector.go
```

SYSTEMD
```
nano /lib/systemd/system/redirector.service  
```

```
[Unit]
Description=redirector
Documentation=https://github.com/almarkov/redirector
After=network.target

[Service]
Type=simple
User=root
ExecStart=/var/projects/redirector/start.sh
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```
service start redirector
```

APACHE
```
a2enmod proxy
a2enmod proxy_http
```

```
nano /etc/apache2/sites-available/000-default.conf
```

```
<VirtualHost *:80>
        ProxyRequests Off
        ProxyPreserveHost On
        ProxyVia Full
        <Proxy *>
                Require all granted
        </Proxy>

   <Location />
      ProxyPass http://127.0.0.1:3060/
      ProxyPassReverse http://127.0.0.1:3060/
   </Location>

        ServerAdmin webmaster@localhost
        DocumentRoot /var/www/html

        ErrorLog ${APACHE_LOG_DIR}/error.log
        CustomLog ${APACHE_LOG_DIR}/access.log combined

</VirtualHost>
```

```
systemctl restart apache2.service
```


Настройка приложения
```
nano /var/projects/redirector/config
```

```
address 127.0.0.1:3060

route /test
PL http://google.pl
TR http://google.ru
DEFAULT http://google.com

route /test2
PL http://yandex.ru
RU http://pornhub.com
DEFAULT http://google.ru
```

```
systemctl restart redirector.service
```
