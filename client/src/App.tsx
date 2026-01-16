import {Container, Stack } from "@chakra-ui/react";
import Navbar from "./components/Navbar";
import TodoForm from "./components/TodoForm";
import TodoList from "./components/TodoList";
export const BASE_URL =
  import.meta.env.Mode === "development" ? "http://localhost:5000/api" : "/api";
// http://localhost:5000 for development and domain+/api for production

function App() {
  return (
    <Stack>
      <Navbar />
      <Container>
        <TodoForm />
        <TodoList />
      </Container>
    </Stack>
  );
}

export default App;
