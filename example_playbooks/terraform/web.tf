resource "aws_instance" "web" {
  ami = "${var.default-ami}"
  instance_type = "t2.small"
  key_name = "colin-linux"
  vpc_security_group_ids = [
    "${var.default-vpc-security-group-id}",
    "${aws_security_group.ssh-access.id}",
    "${aws_security_group.web.id}"
  ]
  subnet_id = "${var.default-subnet-id}"
  private_ip = "172.30.2.11"
  tags = {
    "service" = "test"
  }
}
