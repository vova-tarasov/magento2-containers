# Getting started with an existing project

1. Clone the repository
    ```
    git clone git@github.com:vova-tarasov/magento2-containers.git magento2/
    ```
2. Cd into the newly created folder
    ```shell script
    cd magento2/
    ``` 
3. Clone the project repository
    ```
    git clone git@path-to-git-project.git src/
    ```
3. Build and start the containers 
    ```shell script
   docker-compose up --build 
    ```
4. Connect to the PHP-FPM container
    ```shell script
   docker exec -it $(docker ps -f name=php-fpm -q) bash
    ```
5. Install dependencies
   ```shell script
   composer install
   ```
6. Run installation from a command line
    ```shell script
   bin/magento setup:install \
     --db-host mysql \
     --db-name magento \
     --db-user magento \
     --db-password magento \
     --backend-frontname admin \
     --admin-firstname "admin" \
     --admin-lastname "admin" \
     --admin-email "admin@example.com" \
     --admin-user "admin" \
     --admin-password "123123q" \
     --language "en_US" \
     --currency "USD" \
     --timezone "America/Chicago" \
     --use-rewrites "1" \
     --cleanup-database \
     --http-cache-hosts "varnish" \
     --session-save redis \
     --session-save-redis-host redis \
     --session-save-redis-port 6379 \
     --session-save-redis-db 1 \
     --cache-backend redis \
     --cache-backend-redis-server redis \
     --cache-backend-redis-db 0 \
     --cache-backend-redis-port 6379 \
     --base-url "http://magento2.local/"
    ```
7. Add a new hostname to `/etc/hosts` file

    For MacOS users 
   ```shell script
   echo "127.0.0.1\tmagento2.local" | sudo tee -a /etc/hosts 
   ```

    for Linux users
   ```shell script
   echo -e "127.0.0.1\tmagento2.local" | sudo tee -a /etc/hosts 
   ```

8. Navigate to http://magento2.local

9. Navigate to http://magento2.local/admin to sign in to the backend area
    ```
    Username: admin
    Password: 123123q 
    ```
