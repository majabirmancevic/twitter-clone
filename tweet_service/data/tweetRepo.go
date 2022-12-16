package data

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
)

type TweetRepo struct {
	session *gocql.Session
	logger  *log.Logger
}

func New(logger *log.Logger) (*TweetRepo, error) {
	db := os.Getenv("CASS_DB")

	// Connect to default keyspace
	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Create 'student' keyspace
	err = session.Query(
		fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
					WITH replication = {
						'class' : 'SimpleStrategy',
						'replication_factor' : %d
					}`, "tweet", 1)).Exec()
	if err != nil {
		logger.Println(err)
	}
	session.Close()

	// Connect to student keyspace
	cluster.Keyspace = "tweet"
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &TweetRepo{
		session: session,
		logger:  logger,
	}, nil
}

// Disconnect from database
func (tr *TweetRepo) CloseSession() {
	tr.session.Close()
}

// Create ocene_by_student and ocene_by_predmet tables
func (tr *TweetRepo) CreateTables() {
	err := tr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(regular_username text, description text, id UUID, 
					PRIMARY KEY (regular_username)) `,
			"tweet_by_regular_user")).Exec()
	log.Println("KREIRANJE TABELA")
	if err != nil {
		tr.logger.Println(err)
	}
}

func (tr *TweetRepo) GetTweetsByUser(username string) (TweetsByRegularUser, error) {
	scanner := tr.session.Query(`SELECT regular_username, description, id FROM tweet_by_regular_user WHERE regular_username = ?`,
		username).Iter().Scanner()

	var tweets TweetsByRegularUser
	for scanner.Next() {
		var tweet TweetByRegularUser
		err := scanner.Scan(&tweet.RegularUsername, &tweet.Description, &tweet.Id)
		if err != nil {
			tr.logger.Println(err)
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}
	if err := scanner.Err(); err != nil {
		tr.logger.Println(err)
		return nil, err
	}
	return tweets, nil
}

func (tr *TweetRepo) InsertTweetByRegUser(regUserTweet *TweetByRegularUser) error {
	tweetId, _ := gocql.RandomUUID()
	err := tr.session.Query(
		`INSERT INTO tweet_by_regular_user (regular_username, description, id) 
		VALUES (?, ?, ?)`,
		regUserTweet.RegularUsername, regUserTweet.Description, tweetId).Exec()
	if err != nil {
		tr.logger.Println(err)
		return err
	}
	return nil
}

// NoSQL: Performance issue, we never want to fetch all the data
// (In order to get all student ids we need to contact every partition which are usually located on different servers!)
// Here we are doing it for demonstration purposes (so we can see all student/predmet ids)
func (tr *TweetRepo) GetDistinctIds(idColumnName string, tableName string) ([]string, error) {
	scanner := tr.session.Query(
		fmt.Sprintf(`SELECT DISTINCT %s FROM %s`, idColumnName, tableName)).
		Iter().Scanner()
	var ids []string
	for scanner.Next() {
		var id string
		err := scanner.Scan(&id)
		if err != nil {
			tr.logger.Println(err)
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := scanner.Err(); err != nil {
		tr.logger.Println(err)
		return nil, err
	}
	return ids, nil
}
