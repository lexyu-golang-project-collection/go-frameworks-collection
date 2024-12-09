-- migrate:up
create table artists (
    ArtistId INTEGER not null primary key autoincrement,
    Name NVARCHAR(120)
);
create table albums (
    AlbumId INTEGER not null primary key autoincrement,
    Title NVARCHAR(160) not null,
    ArtistId INTEGER not null references artists
);
create index IFK_AlbumArtistId on albums (ArtistId);
create table employees (
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
create table customers (
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
create index IFK_CustomerSupportRepId on customers (SupportRepId);
create index IFK_EmployeeReportsTo on employees (ReportsTo);
create table genres (
    GenreId INTEGER not null primary key autoincrement,
    Name NVARCHAR(120)
);
create table invoices (
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
create index IFK_InvoiceCustomerId on invoices (CustomerId);
create table media_types (
    MediaTypeId INTEGER not null primary key autoincrement,
    Name NVARCHAR(120)
);
create table playlists (
    PlaylistId INTEGER not null primary key autoincrement,
    Name NVARCHAR(120)
);
create table tracks (
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
create table invoice_items (
    InvoiceLineId INTEGER not null primary key autoincrement,
    InvoiceId INTEGER not null references invoices,
    TrackId INTEGER not null references tracks,
    UnitPrice NUMERIC(10, 2) not null,
    Quantity INTEGER not null
);
create index IFK_InvoiceLineInvoiceId on invoice_items (InvoiceId);
create index IFK_InvoiceLineTrackId on invoice_items (TrackId);
create table playlist_track (
    PlaylistId INTEGER not null references playlists,
    TrackId INTEGER not null references tracks,
    constraint PK_PlaylistTrack primary key (PlaylistId, TrackId)
);
create index IFK_PlaylistTrackTrackId on playlist_track (TrackId);
create index IFK_TrackAlbumId on tracks (AlbumId);
create index IFK_TrackGenreId on tracks (GenreId);
create index IFK_TrackMediaTypeId on tracks (MediaTypeId);
-- migrate:down
drop table invoice_items;
drop table invoices;
drop table customers;
drop table employees;
drop table playlist_track;
drop table playlists;
drop table tracks;
drop table albums;
drop table artists;
drop table genres;
drop table media_types;