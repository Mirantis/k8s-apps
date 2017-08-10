import com.datastax.spark.connector._
import org.apache.spark.SparkContext
import org.apache.spark.SparkContext._
import com.datastax.spark.connector.cql._
import org.apache.spark.SparkConf

case class FoodToUserIndex(food: String, user: String)

val conf = new SparkConf

val c = CassandraConnector(sc.getConf)

val session = c.openSession()

session.execute("CREATE KEYSPACE tutorial WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}")

session.execute("use tutorial")

session.execute("CREATE TABLE tutorial.user (name text primary key, favorite_food text)")

session.execute("create table tutorial.food_to_user_index ( food text, user text, primary key (food, user))")

session.execute("insert into user (name, favorite_food) values ('Jon', 'bacon')")

session.execute("insert into user (name, favorite_food) values ('Luke', 'steak')")

session.execute("insert into user (name, favorite_food) values ('Al', 'salmon')")

session.execute("insert into user (name, favorite_food) values ('Chris', 'chicken')")

session.execute("insert into user (name, favorite_food) values ('Rebecca', 'bacon')")

session.execute("insert into user (name, favorite_food) values ('Patrick', 'brains')")

session.execute("insert into user (name, favorite_food) values ('Duy Hai', 'brains')")

session.close()

val user_table = sc.cassandraTable("tutorial", "user")

user_table.first

val food_index = user_table.map(r => new FoodToUserIndex(r.getString("favorite_food"), r.getString("name")))

food_index.saveToCassandra("tutorial", "food_to_user_index")

val session1 = c.openSession()

session1.execute("drop keyspace tutorial")

session1.close()
