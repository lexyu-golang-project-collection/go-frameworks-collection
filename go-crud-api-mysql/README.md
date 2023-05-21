# Ref
- [GO And MYSQL - 2021 Project 🚀 💣 🔥 - Connect Go with Mysql / Build a Book Management System](https://youtu.be/1E_YycpCsXw)

# Use
- MySQL
- GORM
- Json marshall, unmarshall
- Project structure
- Groilla Mux

```mermaid
flowchart TB
    subgraph PKG
    
    subgraph controllers
    BC[Book-Controller]
    end
    
    subgraph models
    Bg[Book.go]
    end
    
    subgraph routes
    BR[Bookstore-Routes]
    end
    
    subgraph utils
    Ug[Utils.go]
    end

    subgraph config
    Ag[APP.go]
    end
    
    end

    subgraph CMD
    Ag[main.go]
    end
   
    subgraph controller-funcs
    PGB[POST GETBooks `/book/`]
    end
```