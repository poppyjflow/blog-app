import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import styled from 'styled-components';

const PostDetailWrapper = styled.div`
  padding: 2rem;
`;

function PostDetail() {
  const [post, setPost] = useState(null);
  const { id } = useParams();

  useEffect(() => {
    const fetchPost = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/posts/${id}`);
        setPost(response.data);
      } catch (error) {
        console.error('Error fetching post:', error);
      }
    };
    fetchPost();
  }, [id]);

  if (!post) return <div>Loading...</div>;

  return (
    <PostDetailWrapper>
      <h2>{post.title}</h2>
      <p>{post.content}</p>
    </PostDetailWrapper>
  );
}

export default PostDetail;