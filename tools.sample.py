python
import json
import openai # pip install openai requests

client = openai.OpenAI()

# 1. Define the Python functions that act as "tools"
def fetch_weather(city: str) -> str:
    # Simulated or real call to a Weather API (e.g., OpenWeatherMap)
    return json.dumps({"city": city, "temperature": "29°C", "condition": "Sunny"})

def dispatch_email(recipient: str, subject: str, body: str) -> str:
    # Simulated or real call to an Email service (e.g., SendGrid/SMTP/Gmail API)
    return json.dumps({"status": "sent", "to": recipient, "subject": subject})

# 2. Map tool names to actual code functions
available_tools = {
    "fetch_weather": fetch_weather,
    "dispatch_email": dispatch_email,
}

# 3. Describe the tools to the AI model
tools_schema = [
    {
        "type": "function",
        "function": {
            "name": "fetch_weather",
            "description": "Get the current weather for a given city.",
            "parameters": {
                "type": "object",
                "properties": {
                    "city": {"type": "string", "description": "The city name, e.g. Chennai"}
                },
                "required": ["city"]
            }
        }
    },
    {
        "type": "function",
        "function": {
            "name": "dispatch_email",
            "description": "Send an email to a recipient.",
            "parameters": {
                "type": "object",
                "properties": {
                    "recipient": {"type": "string"},
                    "subject": {"type": "string"},
                    "body": {"type": "string"}
                },
                "required": ["recipient", "subject", "body"]
            }
        }
    }
]

def run_ai_agent(user_prompt: string):
    messages = [{"role": "user", "content": user_prompt}]
    
    response = client.chat.completions.create(
        model="gpt-4o-mini",
        messages=messages,
        tools=tools_schema,
        tool_choice="auto" # Let AI choose whether to call a tool
    }
    
    response_message = response.choices[0].message
    tool_calls = response_message.tool_calls

    if tool_calls:
        # AI decided to call a tool
        messages.append(response_message)
        
        for tool_call in tool_calls:
            function_name = tool_call.function.name
            function_args = json.loads(tool_call.function.arguments)
            
            # Execute the matching Python function
            if function_name in available_tools:
                tool_to_call = available_tools[function_name]
                tool_result = tool_to_call(**function_args)
                
                messages.append({
                    "role": "tool",
                    "tool_call_id": tool_call.id,
                    "name": function_name,
                    "content": tool_result
                })
        
        # Second call to AI to synthesize the tool result into a natural response
        second_response = client.chat.completions.create(
            model="gpt-4o-mini",
            messages=messages
        )
        return second_response.choices[0].message.content
    
    return response_message.content

# --- Scenario A: Triggers Weather API ---
print(run_ai_agent("What is the weather like in Chennai right now?"))
# AI calls `fetch_weather(city='Chennai')`

# --- Scenario B: Triggers Email Tool ---
print(run_ai_agent("Send an email to alice@example.com with the subject 'Meeting' and body 'See you at 4 PM'"))
# AI calls `dispatch_email(recipient='alice@example.com', ...)`