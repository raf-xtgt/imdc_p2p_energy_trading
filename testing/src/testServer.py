from flask import Flask

server_port = 5000
app = Flask(__name__)

@app.route('/login/<username>/<password>')
def login(username, password):
    return f'Hello, {username}:{password}!'

@app.route('/logout')
def logout():
  return f'Bye {username}:{password}!'

@app.route('/')
def index():
  return 'Index'

@app.route('/test')
def test():
  return 'Test'
if __name__ == "__main__":
    app.run('0.0.0.0',port=server_port)