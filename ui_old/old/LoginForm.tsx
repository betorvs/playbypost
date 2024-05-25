import React, { useState } from "react";

interface LoginFormProps {
  setToken: (
    access_token: string,
    expire_on: EpochTimeStamp,
    user_id: number,
    username: string
  ) => void;
}

interface Session {
  access_token: string;
  expire_on: EpochTimeStamp;
  user_id: number;
  status: string;
  msg: string;
}

const LoginForm = ({ setToken }: LoginFormProps) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    console.log(event);

    // Implement login logic here (e.g., API call)
    // Replace with your actual authentication logic
    try {
      const urlAPI = new URL("login", "http://192.168.1.210:3000");
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
      //   console.log(response.json());

      if (response.ok) {
        // Login successful, handle redirection or state change
        console.log("Login successful!");
        const data = await response.text();
        const res: Session = JSON.parse(data);
        console.log(response);
        setToken(res.access_token, res.expire_on, res.user_id, username);
      } else {
        setErrorMessage("Invalid username or password");
      }
    } catch (error) {
      console.error("Login error:", error);
      setErrorMessage("An error occurred. Please try again.");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="mb-3">
        <label htmlFor="text" className="form-label">
          Username
        </label>
        <input
          type="uname"
          className="form-control"
          id="username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
      </div>
      <div className="mb-3">
        <label htmlFor="password" className="form-label">
          Password
        </label>
        <input
          type="password"
          className="form-control"
          id="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </div>
      <button type="submit" className="btn btn-primary">
        Login
      </button>
      {errorMessage && (
        <div className="alert alert-danger mt-3">{errorMessage}</div>
      )}
    </form>
  );
};

export default LoginForm;
