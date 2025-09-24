from faker import Faker
import random

fake = Faker()

for _ in range(1000):
    name = fake.first_name()
    email = fake.email()
    phone = fake.msisdn()[:10]
    attending = random.choice(["yes", "no"])
    
    print(name)
    print(email)
    print(phone)
    print(attending)

print("exit")
print("Chris")
