# Course System

This is the project of Group 28 for ByteCamp 2022 Winter.

## Quick Start

1. Clone the project.
    
    ```
    git clone https://github.com/Dexter1636/course_system.git
    ```

2. Download modules.
    
    ```
    cd course_system
    go mod download
    ```

3. Add application config file and test config file.
    
    Write the following code to your `course_system/config/application.yaml`:
    
    ```yaml
    server:
      port: 8080
    
    datasource:
      driverName: mysql
      host: <hostname>
      port: <port>
      database: <database_name>
      username: <username>
      password: <password>
      charset: utf8
   
   redis:
      host: <hostname>
      port: <port>
      db: <db>
      user: <username>
      password: <password>
    
    logger:
      level: info
    ```
   
    And the same for `course_system/config/test.yaml`, which is used for testing.

4. Run.
    
    ```
    go run .
    ```

## Note

1. Remember to use your own branch for development.
    
    Do this before coding:
    
    ```
    git checkout dev
    git branch <your_branch_name>
    git checkout <your_branch_name>
    git push --set-upstream origin <your_branch_name>
    ```
   
2. Do NOT track `application.yaml` and `test.yaml` since they contain sensitive data.
