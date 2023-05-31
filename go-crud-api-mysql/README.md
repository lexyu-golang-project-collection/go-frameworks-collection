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
    post[POST]
    get[GET]
    put[PUT]
    delete[DELETE]

    CBs[Create Books]
    GBs[Get Books]
    GBBI[Get Book By ID]
    UB[Update Book]
    DB[Delete Book]

    PCBs[`/book/`]
    GGBs[`/book/`]
    GGBBI[`/book/bookId`]
    PUB[`/book/bookId`]
    DDB[`/book/bookId`]

    post --> CBs --> PCBs
    get --> GBs --> GGBs
    get --> GBBI --> GGBBI
    put --> UB --> PUB
    delete --> DB --> DDB

    end
```