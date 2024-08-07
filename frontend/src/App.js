import React from 'react';
import { BrowserRouter as Router, Routes, Route, Switch } from 'react-router-dom';
import styled, { ThemeProvider } from 'styled-components';
import { darkTheme } from './theme';
import Header from './components/Header';
import Home from './components/Home';
import CreatePost from './components/CreatePost';
import PostList from './components/PostList';
import PostDetail from './components/PostDetail';

const AppWrapper = styled.div`
  background-color: ${props => props.theme.background};
  color: ${props => props.theme.text};
  min-height: 100vh;
`;

function App() {
  return (
    <ThemeProvider theme={darkTheme}>
      <AppWrapper>
        <Router>
          <Header />
          <Routes>
            <Route exact path="/" element={<Home />} />
            <Route path="/create" element={<CreatePost />} />
            <Route exact path="/posts" element={<PostList />} />
            <Route path="/posts/:id" element={<PostDetail />} />
          </Routes>
        </Router>
      </AppWrapper>
    </ThemeProvider>
  );
}

export default App;
