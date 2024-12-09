CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE artists (
    ArtistId INTEGER not null primary key autoincrement,
    Name NVARCHAR(120)
);
CREATE TABLE albums (
    AlbumId INTEGER not null primary key autoincrement,
    Title NVARCHAR(160) not null,
    ArtistId INTEGER not null references artists
);
CREATE INDEX IFK_AlbumArtistId on albums (ArtistId);
CREATE TABLE employees (
    EmployeeId INTEGER not null primary key autoincrement,
    LastName NVARCHAR(20) not null,
    FirstName NVARCHAR(20) not null,
    Title NVARCHAR(30),
    ReportsTo INTEGER references employees,
    BirthDate DATETIME,
    HireDate DATETIME,
    Address NVARCHAR(70),
    City NVARCHAR(40),
    State NVARCHAR(40),
    Country NVARCHAR(40),
    PostalCode NVARCHAR(10),
    Phone NVARCHAR(24),
    Fax NVARCHAR(24),
    Email NVARCHAR(60)
);
CREATE TABLE customers (
    CustomerId INTEGER not null primary key autoincrement,
    FirstName NVARCHAR(40) not null,
    LastName NVARCHAR(20) not null,
    Company NVARCHAR(80),
    Address NVARCHAR(70),
    City NVARCHAR(40),
    State NVARCHAR(40),
    Country NVARCHAR(40),
    PostalCode NVARCHAR(10),
    Phone NVARCHAR(24),
    Fax NVARCHAR(24),
    Email NVARCHAR(60) not null,
    SupportRepId INTEGER references employees
);
CREATE INDEX IFK_CustomerSupportRepId on customers (SupportRepId);
CREATE INDEX IFK_EmployeeReportsTo on employees (ReportsTo);
CREATE TABLE genres (
    GenreId INTEGER not null primary key autoincrement,
    Name NVARCHAR(120)
);
CREATE TABLE invoices (
    InvoiceId INTEGER not null primary key autoincrement,
    CustomerId INTEGER not null references customers,
    InvoiceDate DATETIME not null,
    BillingAddress NVARCHAR(70),
    BillingCity NVARCHAR(40),
    BillingState NVARCHAR(40),
    BillingCountry NVARCHAR(40),
    BillingPostalCode NVARCHAR(10),
    Total NUMERIC(10, 2) not null
);
CREATE INDEX IFK_InvoiceCustomerId on invoices (CustomerId);
CREATE TABLE media_types (
    MediaTypeId INTEGER not null primary key autoincrement,
    Name NVARCHAR(120)
);
CREATE TABLE playlists (
    PlaylistId INTEGER not null primary key autoincrement,
    Name NVARCHAR(120)
);
CREATE TABLE tracks (
    TrackId INTEGER not null primary key autoincrement,
    Name NVARCHAR(200) not null,
    AlbumId INTEGER references albums,
    MediaTypeId INTEGER not null references media_types,
    GenreId INTEGER references genres,
    Composer NVARCHAR(220),
    Milliseconds INTEGER not null,
    Bytes INTEGER,
    UnitPrice NUMERIC(10, 2) not null
);
CREATE TABLE invoice_items (
    InvoiceLineId INTEGER not null primary key autoincrement,
    InvoiceId INTEGER not null references invoices,
    TrackId INTEGER not null references tracks,
    UnitPrice NUMERIC(10, 2) not null,
    Quantity INTEGER not null
);
CREATE INDEX IFK_InvoiceLineInvoiceId on invoice_items (InvoiceId);
CREATE INDEX IFK_InvoiceLineTrackId on invoice_items (TrackId);
CREATE TABLE playlist_track (
    PlaylistId INTEGER not null references playlists,
    TrackId INTEGER not null references tracks,
    constraint PK_PlaylistTrack primary key (PlaylistId, TrackId)
);
CREATE INDEX IFK_PlaylistTrackTrackId on playlist_track (TrackId);
CREATE INDEX IFK_TrackAlbumId on tracks (AlbumId);
CREATE INDEX IFK_TrackGenreId on tracks (GenreId);
CREATE INDEX IFK_TrackMediaTypeId on tracks (MediaTypeId);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20241126130434');
