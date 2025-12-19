resource "aws_db_instance" "postgres" {
  allocated_storage    = 20
  db_name              = "appdb"
  engine               = "postgres"
  engine_version       = "15.4"
  instance_class       = "db.t3.micro"
  username             = "dbadmin"
  password             = "securepassword" # In real usage, use Secrets Manager or Variables
  parameter_group_name = "default.postgres15"
  skip_final_snapshot  = true
  publicly_accessible  = false
}
