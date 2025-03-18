upstream stockub {
        server localhost:8081;
}

server {
        server_name ecg-online.suttapak.com;

        location / {
                client_max_body_size 24M;
                proxy_pass http://stockub;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Host $host;
                proxy_cache_bypass $http_upgrade;
                proxy_set_header Connection �~@~Xupgrade�~@~Y;
        }

}


