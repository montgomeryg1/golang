package pedapilog

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)


type Application struct {
    DB *pgxpool.Pool
}


func (app *Application) WorldMap(ctx context.Context) (error) {
    tx, err := app.DB.Begin(ctx)
    if err != nil {
        return err
    }

    // Defining duration 
    d := (60 * time.Minute) 

    // Calling Round() method 
    then := time.Now().UTC().Add(-2 * time.Hour).Round(d).Format("2006-01-02 15:04")
    now := time.Now().UTC().Round(d).Format("2006-01-02 15:04")

    qry := fmt.Sprintf("INSERT INTO worldmap SELECT geoloc.latitude, geoloc.longitude, geoloc.city, httplogs.requesttime, httplogs.ipaddress FROM httplogs LEFT JOIN geoloc ON httplogs.ipaddress = geoloc.ipaddress WHERE httplogs.requesttime between '%s' and '%s'", then, now)
    qry += " ON CONFLICT ON CONSTRAINT ux_ip_req DO NOTHING"
    _, err = tx.Exec(ctx, qry)
    if err != nil {
        tx.Rollback(ctx)
        return err
    } 
    
    err = tx.Commit(ctx)
    return err
}


func (app *Application) DeleteOld(ctx context.Context) (error) {
    // // Defining duration 
    // d := (24 * time.Hour) 

    // // Calling Round() method 
    // httplogExp := time.Now().UTC().Add(-168 * time.Hour).Round(d).Format("2006-01-02")
    // worldmapExp := time.Now().UTC().Add(-48 * time.Hour).Round(d).Format("2006-01-02")
    tx, err := app.DB.Begin(ctx)
    if err != nil {
        return err
    }

    // qry := fmt.Sprintf("DELETE FROM httplogs WHERE requesttime < '%s'", httplogExp)
    qry := "DELETE FROM httplogs where requesttime < NOW() - INTERVAL '8 day'"
    _, err = tx.Exec(ctx, qry)
    if err != nil {
        tx.Rollback(ctx)
        return err
    } 
    
    //qry = fmt.Sprintf("DELETE FROM worldmap WHERE requesttime < '%s'",worldmapExp)
    qry = "DELETE FROM worldmap where requesttime < NOW() - INTERVAL '3 day'"
    _, err = tx.Exec(ctx, qry)
    if err != nil {
        tx.Rollback(ctx)
        return err
    } 

    err = tx.Commit(ctx)
    return err
}

// Requests by PartnerID
func (app *Application) RequestsByPartnerID(ctx context.Context) (error) {
    tx, err := app.DB.Begin(ctx)
    if err != nil {
        return err
    }

    // Defining duration 
    d := (24 * time.Hour) 

    // Calling Round() method 
    then := time.Now().UTC().Add(-24 * time.Hour).Round(d).Format("2006-01-02")
    now := time.Now().UTC().Round(d).Format("2006-01-02")

    qry := fmt.Sprintf("insert into requests_by_partnerid (select date_trunc('day', requesttime), partnerid, count(1) from httplogs where requesttime between '%s' and '%s'", then, now)
    qry += " group by date_trunc, partnerid) ON CONFLICT ON CONSTRAINT unq_date_partnerid DO UPDATE SET count = excluded.count"
    _, err = tx.Exec(ctx, qry)
    if err != nil {
        tx.Rollback(ctx)
        return err
    } 
    
    err = tx.Commit(ctx)
    return err
}