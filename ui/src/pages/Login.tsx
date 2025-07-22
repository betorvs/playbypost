import React, { useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../context/AuthContext";
import { Form } from "react-bootstrap";
import Button from "react-bootstrap/Button";
import SaveToken from "../context/SaveToken";
import UseLocation from "../context/UseLocation";
import { SessionToken } from "../types/Session";

const Login = () => {
  const { setAuthenticated } = useContext(AuthContext);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const navigate = useNavigate();

  const handleLogin = () => {
    setAuthenticated(true);
    navigate("/");
  };

  async function clickLogin(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/login", apiURL);
    console.log("urlAPI to call login", urlAPI);
    const response = await fetch(urlAPI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: username,
        password: password,
      }),
    });
    if (response.ok) {
      alert("You are logged in");
      const data = await response.text();
      const res: SessionToken = JSON.parse(data);
      SaveToken(res.access_token, res.expire_on, res.user_id, username);
      handleLogin();
    } else {
      alert("Please check your login information.");
    }
  }

  return (
    <div className="container mt-3" key="1">
      <h1>Play by Post Login</h1>
      <hr />
      <Form onSubmit={clickLogin}>
        <Form.Group className="mb-3" controlId="formUsername">
          <Form.Label>Username</Form.Label>
          <Form.Control
            type="text"
            placeholder="Enter username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <Form.Text className="text-muted">
            We'll never share your username with anyone else.
          </Form.Text>
        </Form.Group>

        <Form.Group className="mb-3" controlId="formPassword">
          <Form.Label>Password</Form.Label>
          <Form.Control
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </Form.Group>
        <Button variant="primary" type="submit">
          Submit
        </Button>
      </Form>
    </div>
  );
};

export default Login;
