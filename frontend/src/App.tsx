import { useEffect, useState } from 'react'
import './App.css'
import axios, { AxiosError } from 'axios'

function App() {
  const [message, setMessage] = useState("")
  const [username, setUsername] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  function onSubmit() {
    axios.post("http://localhost:8080/user/register", {
      username: username,
      email: email,
      password: password
    }).then((response) => {
      console.log(response.data.message);
      console.log(response.data.data);
      let message: string = response.data.message;
      setMessage(message)
    }).catch((error: AxiosError) => {
      alert(error)
    })
  }

  function onSubmitLogin() {
    axios.post("http://localhost:8080/user/login", {
      username: username,
      email: email,
      password: password
    }).then((response) => {
      console.log(response.data.message);
      console.log(response.data.data);
      let message: string = response.data.message;
      setMessage(message)
    }).catch((error: AxiosError) => {
      alert(error)
    })
  }
  

  return (
    <>
      <h1>Create User</h1>
      
      <form action="" onSubmit={onSubmit}>
        <input type="text" name="username" id="username" placeholder="Tom" value={username} onChange={(e) => setUsername(e.target.value)}/><br/>
        <input type="text" name="email" id="email" placeholder="Tom@gmail.com" value={email} onChange={(e) => setEmail(e.target.value)}/><br/>
        <input type="text" name="password" id="password" placeholder="123456743" value={password} onChange={(e) => setPassword(e.target.value)}/><br/>
        <input type="submit" value="Create Account" /><br/>
      </form>
      {message}
      
      <br/>
      <form action="" onSubmit={onSubmitLogin}>
        <input type="text" name="username" id="username" placeholder="Tom" value={username} onChange={(e) => setUsername(e.target.value)}/><br/>
        <input type="email" name="email" id="email" placeholder="Tom@gmail.com" value={email} onChange={(e) => setEmail(e.target.value)}/><br/>
        <input type="text" name="password" id="password" placeholder="123456743" value={password} onChange={(e) => setPassword(e.target.value)}/><br/>
        <input type="submit" value="Login" /><br/>
      </form>
      {message}
    </>
  )
}

export default App
