#!/usr/bin/env python3
"""
Automatically fill your RSVP form with fake data using Python Faker
"""

import requests
from faker import Faker
import random
import time

# Initialize Faker
fake = Faker()

# Your Railway app URL
BASE_URL = "https://hands-on-golang-production.up.railway.app"

def generate_rsvp_data():
    """Generate realistic RSVP data"""
    return {
        'name': fake.name(),
        'email': fake.email(),
        'phone': fake.phone_number(),
        'willAttend': random.choice(['true', 'false'])
    }

def submit_rsvp(rsvp_data):
    """Submit RSVP data to your form"""
    try:
        response = requests.post(
            f"{BASE_URL}/form",
            data=rsvp_data,
            headers={
                'Content-Type': 'application/x-www-form-urlencoded',
                'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
            }
        )
        
        if response.status_code == 200:
            status = "✅ ATTENDING" if rsvp_data['willAttend'] == 'true' else "❌ NOT ATTENDING"
            print(f"{status} - {rsvp_data['name']} ({rsvp_data['email']})")
            return True
        else:
            print(f"❌ Failed to submit {rsvp_data['name']}: HTTP {response.status_code}")
            return False
            
    except requests.RequestException as e:
        print(f"❌ Error submitting {rsvp_data['name']}: {e}")
        return False

def bulk_populate_rsvps(count=10):
    """Generate and submit multiple RSVPs"""
    print(f"🚀 Starting bulk RSVP generation for {count} people...")
    print(f"🎯 Target: {BASE_URL}")
    print("-" * 60)
    
    successful = 0
    failed = 0
    
    for i in range(count):
        rsvp_data = generate_rsvp_data()
        
        if submit_rsvp(rsvp_data):
            successful += 1
        else:
            failed += 1
            
        # Small delay to be nice to your server
        time.sleep(0.5)
    
    print("-" * 60)
    print(f"📊 Results: {successful} successful, {failed} failed")
    print(f"🎉 Check your guest list at: {BASE_URL}/list")

def interactive_mode():
    """Interactive mode for custom control"""
    print("🎭 Interactive RSVP Generator")
    print("Commands:")
    print("  - Number (e.g., '5'): Generate that many RSVPs")
    print("  - 'single': Generate one RSVP")
    print("  - 'list': View current guest list")
    print("  - 'quit': Exit")
    
    while True:
        command = input("\n> ").strip().lower()
        
        if command == 'quit':
            break
        elif command == 'single':
            rsvp_data = generate_rsvp_data()
            submit_rsvp(rsvp_data)
        elif command == 'list':
            print(f"🔗 Guest list: {BASE_URL}/list")
        elif command.isdigit():
            count = int(command)
            if 1 <= count <= 50:  # Reasonable limits
                bulk_populate_rsvps(count)
            else:
                print("❌ Please enter a number between 1 and 50")
        else:
            print("❌ Unknown command")

def demo_mixed_responses():
    """Generate a realistic mix of attending vs not attending"""
    print("🎪 Generating realistic party responses...")
    
    # Generate more "yes" responses (70% attending is realistic for a party)
    attendees = []
    non_attendees = []
    
    for i in range(15):  # 15 people total
        rsvp_data = generate_rsvp_data()
        
        # 70% chance of attending
        if random.random() < 0.7:
            rsvp_data['willAttend'] = 'true'
            attendees.append(rsvp_data)
        else:
            rsvp_data['willAttend'] = 'false'
            non_attendees.append(rsvp_data)
    
    print(f"📋 Generated: {len(attendees)} attending, {len(non_attendees)} not attending")
    
    # Submit all responses
    all_responses = attendees + non_attendees
    random.shuffle(all_responses)  # Random submission order
    
    for rsvp in all_responses:
        submit_rsvp(rsvp)
        time.sleep(0.3)

if __name__ == "__main__":
    print("🎉 Welcome to the RSVP Auto-Filler!")
    print("Choose a mode:")
    print("1. Quick Demo (15 realistic responses)")
    print("2. Bulk Generate (specify count)")
    print("3. Interactive Mode")
    
    choice = input("\nEnter choice (1-3): ").strip()
    
    if choice == "1":
        demo_mixed_responses()
    elif choice == "2":
        try:
            count = int(input("How many RSVPs to generate? "))
            bulk_populate_rsvps(count)
        except ValueError:
            print("❌ Please enter a valid number")
    elif choice == "3":
        interactive_mode()
    else:
        print("❌ Invalid choice, running quick demo...")
        demo_mixed_responses()
    
    print(f"\n🎊 All done! Check your results at: {BASE_URL}/list")