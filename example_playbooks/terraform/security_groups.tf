resource "aws_security_group" "ssh-access" {
  name = "dev ssh access"
  description = "allow access to ssh"
  vpc_id = "${ var.vpc-id }"

  ingress {
    from_port = 22
    to_port = 22
    protocol = "TCP"
    cidr_blocks = ["${var.home_ip_address}"]
  }

  tags {
    Name = "ssh-access"
  }
}

resource "aws_security_group" "web" {
  name = "web"
  description = "web security group"
  vpc_id = "${ var.vpc-id }"

  ingress {
    from_port = 443
    to_port = 443
    protocol = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 80
    to_port = 80
    protocol = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
