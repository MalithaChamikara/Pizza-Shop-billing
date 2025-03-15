import {
  Container,
  Paper,
  Typography,
  Box,
  TextField,
  Button,
} from "@mui/material";

const LoginPage = () => {
  const handleSubmit = () => console.log("login");
  return (
    <Container maxWidth="xs">
      <Paper elevation={10} sx={{ marginTop: 8, padding: 2 }}>
        <Typography component="h1" variant="h5" sx={{ textAlign: "center" }}>
          LOG IN
        </Typography>
        <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
          <TextField
            placeholder="username"
            fullWidth
            required
            autoFocus
            sx={{ mb: 2 }}
          ></TextField>
          <TextField
            placeholder="password"
            fullWidth
            required
            type="password"
            sx={{ mb: 2 }}
          ></TextField>
          <Button type="submit" variant="contained" fullWidth sx={{ mt: 1 }}>
            Login
          </Button>
        </Box>
      </Paper>
    </Container>
  );
};

export default LoginPage;
