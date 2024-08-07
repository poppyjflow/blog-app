import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';
import styled from 'styled-components';

const PostListWrapper = styled.div`
  padding: 2rem;
`;

const PostItem = styled.div`
  margin-bottom: 1rem;
  padding: 1rem;
  background-color: ${props => props.theme.secondary};
  border-radius: 4px;
`;

function PostList() {
  const [posts, setPosts] = useState([]);

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const response = await axios.get('http://localhost:8080/posts');
        setPosts(response.data);
      } catch (error) {
        console.error('Error fetching posts:', error);
      }
    };
    fetchPosts();
  }, []);

  return (
    <PostListWrapper>
      <h2>All Posts</h2>
      {posts.map(post => (
        <PostItem key={post.id}>
          <h3>{post.title}</h3>
          <Link to={`/posts/${post.id}`}>Read more</Link>
        </PostItem>
      ))}
    </PostListWrapper>
  );
}

export default PostList;