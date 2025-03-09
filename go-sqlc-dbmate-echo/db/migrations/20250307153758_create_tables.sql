-- migrate:up
CREATE TABLE pages (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    slug         VARCHAR(255) NOT NULL UNIQUE,
    status       ENUM ('draft', 'published', 'archived') NOT NULL DEFAULT 'draft',
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL,

    INDEX idx_page_status (status),
    INDEX idx_page_slug (slug)
);

CREATE TABLE page_elements (
    id               INT AUTO_INCREMENT PRIMARY KEY,
    uuid             VARCHAR(36) NOT NULL UNIQUE,
    page_id          INT         NOT NULL,
    block_type     ENUM ('h1', 'h2', 'h3', 'h4', 'h5', 'h6', 
                           'paragraph', 'ul', 'code', 'quote', 'image') NOT NULL,
    order_index      INT         NOT NULL AUTO_INCREMENT, -- 只用 order 控制順序
    html_content     TEXT    NULL,  -- Markdown 轉成 HTML 的內容
    attributes       JSON        NULL,  -- 圖片URL, 語言類型 (程式碼區塊)
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (page_id) REFERENCES pages (id) ON DELETE CASCADE,

    INDEX idx_uuid (uuid),
    INDEX idx_page_order (page_id, order_index),
    UNIQUE KEY unique_page_order (page_id, order_index)
);


-- migrate:down
DROP TABLE IF EXISTS page_elements;
DROP TABLE IF EXISTS page_contents;
DROP TABLE IF EXISTS pages;

