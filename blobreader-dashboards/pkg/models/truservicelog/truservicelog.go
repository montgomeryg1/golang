package truservicelog

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// TruserviceRequestModel type which wraps a sql.DB connection pool.
type Application struct {
    DB *pgxpool.Pool
}

// LastSeen will update the lastseen table.
func (app *Application) MerchantLastSeen(ctx context.Context) (error) {
    tx, err := app.DB.Begin(ctx)
    if err != nil {
        return err
    }

    _, err = tx.Exec(ctx,"DELETE FROM lastseen")
    if err != nil {
        // If there is any error, we call the tx.Rollback() method on the
        // transaction. This will abort the transaction and no changes will be
        // made to the database.
        tx.Rollback(ctx)
        return err
    }
    
    _, err = tx.Exec(ctx,"INSERT INTO lastseen select t.id, t.partnerid, t.merchantid, t.requesttime, t.messagetype, t.infomessage, DATE_PART('day', now()::timestamp - requesttime::timestamp) AS dayssincelastseen from httplogs t inner join (select partnerid, merchantid, max(requesttime) as MaxDate from httplogs group by partnerid, merchantid ) tm on t.partnerid = tm.partnerid and t.merchantid = tm.merchantid and t.requesttime = tm.MaxDate ORDER BY requesttime asc")
    if err != nil {
        tx.Rollback(ctx)
        return err
    }
    
    err = tx.Commit(ctx)
    return err
}



// GetRatingsByHour
func (app *Application) GetRatingsByHour(ctx context.Context) (error) {
    // Defining duration 
    d := (60 * time.Minute) 

    // Calling Round() method 
    then := time.Now().UTC().Add(-2 * time.Hour).Round(d).Format("2006-01-02 15:04")
    now := time.Now().UTC().Round(d).Format("2006-01-02 15:04")

    tx, err := app.DB.Begin(ctx)
    if err != nil {
        return err
    }
    
    qry := fmt.Sprintf("INSERT INTO customer_ratings_by_hour select date_trunc('hour', requesttime), count(1) from httplogs where requesttime between '%s' and '%s'", then, now)
    qry += " and messagetype = 'Rating' and requestbody NOT LIKE '%\"-%' group by 1 ON CONFLICT (date_trunc) DO UPDATE SET count = excluded.count"
    _, err = tx.Exec(ctx,qry)
    // _, err = tx.Exec("INSERT INTO customer_ratings_by_hour select date_trunc('hour', requesttime), count(1) from httplogs where requesttime >= NOW() - '2 hour'::INTERVAL and messagetype = 'Rating' and requestbody NOT LIKE '%\"-%' group by 1 ON CONFLICT (date_trunc) DO UPDATE SET count = excluded.count")
    if err != nil {
        tx.Rollback(ctx)
        return err
    }
    
    qry = fmt.Sprintf("INSERT INTO customer_skips_by_hour select date_trunc('hour', requesttime), count(1) from httplogs where requesttime between '%s' and '%s'", then, now)
    qry += " and messagetype = 'Rating' and requestbody LIKE '%\"-1%' group by 1 ON CONFLICT (date_trunc) DO UPDATE SET count = excluded.count"
    _, err = tx.Exec(ctx,qry)
    // _, err = tx.Exec("INSERT INTO customer_skips_by_hour select date_trunc('hour', requesttime), count(1) from httplogs where requesttime >= NOW() - '2 hour'::INTERVAL and messagetype = 'Rating' and requestbody LIKE '%\"-1%' group by 1 ON CONFLICT (date_trunc) DO UPDATE SET count = excluded.count")
    if err != nil {
        tx.Rollback(ctx)
        return err
    } 

    qry = fmt.Sprintf("INSERT INTO question_timeouts_by_hour select date_trunc('hour', requesttime), count(1) from httplogs where requesttime between '%s' and '%s'", then, now)
    qry += " and messagetype = 'Rating' and requestbody LIKE '%\"-2%' group by 1 ON CONFLICT (date_trunc) DO UPDATE SET count = excluded.count"
    _, err = tx.Exec(ctx,qry)
    // // _, err = tx.Exec("INSERT INTO question_timeouts_by_hour select date_trunc('hour', requesttime), count(1) from httplogs where requesttime >= NOW() - '2 hour'::INTERVAL and messagetype = 'Rating' and requestbody LIKE '%\"-2%' group by 1  ON CONFLICT (date_trunc) DO UPDATE SET count = excluded.count")
    if err != nil {
        tx.Rollback(ctx)
        return err
    }        
    
    err = tx.Commit(ctx)
    return err
}

// WorldMap
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

    qry := fmt.Sprintf("INSERT INTO worldmap SELECT httplogs.ipaddress, geoloc.latitude, geoloc.longitude, geoloc.city, httplogs.requesttime FROM httplogs LEFT JOIN geoloc ON httplogs.ipaddress = geoloc.ipaddress WHERE httplogs.requesttime between '%s' and '%s'", then, now)
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

    //qry := fmt.Sprintf("DELETE FROM httplogs WHERE requesttime < '%s'", httplogExp)
    qry := "DELETE FROM httplogs where requesttime < NOW() - INTERVAL '8 day'"
    _, err = tx.Exec(ctx, qry)
    if err != nil {
        tx.Rollback(ctx)
        return err
    } 
    
    // qry = fmt.Sprintf("DELETE FROM worldmap WHERE requesttime < '%s'",worldmapExp)
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