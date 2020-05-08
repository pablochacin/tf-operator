variable "greetee" {
  type = string
}

output "greetings" {

    value = "Hello ${var.greetee}"

}
