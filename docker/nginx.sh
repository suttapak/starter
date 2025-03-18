
nginx -t
sudo certbot --nginx -d stockub.suttapak.com
 systemctl restart nginx
 sudo systemctr restart nginx
 sudo ln -s /etc/nginx/sites-available/service-insider.labotronmedical.com /etc/nginx/sites-enabled/service-insider.labotronmedical.com