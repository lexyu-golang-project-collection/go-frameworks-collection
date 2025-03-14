// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlite

import (
	"database/sql"
	"time"
)

type Album struct {
	Albumid  int64  `json:"albumid"`
	Title    string `json:"title"`
	Artistid int64  `json:"artistid"`
}

type Artist struct {
	Artistid int64          `json:"artistid"`
	Name     sql.NullString `json:"name"`
}

type Customer struct {
	Customerid   int64          `json:"customerid"`
	Firstname    string         `json:"firstname"`
	Lastname     string         `json:"lastname"`
	Company      sql.NullString `json:"company"`
	Address      sql.NullString `json:"address"`
	City         sql.NullString `json:"city"`
	State        sql.NullString `json:"state"`
	Country      sql.NullString `json:"country"`
	Postalcode   sql.NullString `json:"postalcode"`
	Phone        sql.NullString `json:"phone"`
	Fax          sql.NullString `json:"fax"`
	Email        string         `json:"email"`
	Supportrepid sql.NullInt64  `json:"supportrepid"`
}

type Employee struct {
	Employeeid int64          `json:"employeeid"`
	Lastname   string         `json:"lastname"`
	Firstname  string         `json:"firstname"`
	Title      sql.NullString `json:"title"`
	Reportsto  sql.NullInt64  `json:"reportsto"`
	Birthdate  sql.NullTime   `json:"birthdate"`
	Hiredate   sql.NullTime   `json:"hiredate"`
	Address    sql.NullString `json:"address"`
	City       sql.NullString `json:"city"`
	State      sql.NullString `json:"state"`
	Country    sql.NullString `json:"country"`
	Postalcode sql.NullString `json:"postalcode"`
	Phone      sql.NullString `json:"phone"`
	Fax        sql.NullString `json:"fax"`
	Email      sql.NullString `json:"email"`
}

type Genre struct {
	Genreid int64          `json:"genreid"`
	Name    sql.NullString `json:"name"`
}

type Invoice struct {
	Invoiceid         int64          `json:"invoiceid"`
	Customerid        int64          `json:"customerid"`
	Invoicedate       time.Time      `json:"invoicedate"`
	Billingaddress    sql.NullString `json:"billingaddress"`
	Billingcity       sql.NullString `json:"billingcity"`
	Billingstate      sql.NullString `json:"billingstate"`
	Billingcountry    sql.NullString `json:"billingcountry"`
	Billingpostalcode sql.NullString `json:"billingpostalcode"`
	Total             interface{}    `json:"total"`
}

type InvoiceItem struct {
	Invoicelineid int64       `json:"invoicelineid"`
	Invoiceid     int64       `json:"invoiceid"`
	Trackid       int64       `json:"trackid"`
	Unitprice     interface{} `json:"unitprice"`
	Quantity      int64       `json:"quantity"`
}

type MediaType struct {
	Mediatypeid int64          `json:"mediatypeid"`
	Name        sql.NullString `json:"name"`
}

type Playlist struct {
	Playlistid int64          `json:"playlistid"`
	Name       sql.NullString `json:"name"`
}

type PlaylistTrack struct {
	Playlistid int64 `json:"playlistid"`
	Trackid    int64 `json:"trackid"`
}

type SchemaMigrations struct {
	Version string `json:"version"`
}

type Track struct {
	Trackid      int64          `json:"trackid"`
	Name         string         `json:"name"`
	Albumid      sql.NullInt64  `json:"albumid"`
	Mediatypeid  int64          `json:"mediatypeid"`
	Genreid      sql.NullInt64  `json:"genreid"`
	Composer     sql.NullString `json:"composer"`
	Milliseconds int64          `json:"milliseconds"`
	Bytes        sql.NullInt64  `json:"bytes"`
	Unitprice    interface{}    `json:"unitprice"`
}
