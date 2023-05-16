import {
  Button,
  Input,
  PressEvent,
  Spacer,
  Table,
  useInput,
} from '@nextui-org/react';
import { useState } from 'react';

interface BlogPost {
  title: string;
  url: string;
}

export default function Home() {
  const blogPostsColumns = [
    {
      key: 'title',
      label: 'Title',
    },
    {
      key: 'url',
      label: 'URL',
    },
  ];

  const blogPostsTableFC = (label: string, state: BlogPost[]) => {
    return (
      <Table
        aria-label={label}
      >
        <Table.Header columns={blogPostsColumns}>
          {(column) => (
            <Table.Column key={column.key}>{column.label}</Table.Column>
          )}
        </Table.Header>
        <Table.Body items={state}>
          {(item: BlogPost) => (
            <Table.Row key={item.url}>
              <Table.Cell>{item.title}</Table.Cell>
              <Table.Cell><a href={item.url}>{item.url}</a></Table.Cell>
            </Table.Row>
          )}
        </Table.Body>
      </Table>
    );
  };

  const [scrapedBlogPosts, setScrapedBlogPosts] = useState([]);

  const scrapeButtonOnPress = async (_: PressEvent) => {
    const resp = await fetch('http://localhost:8000/api/scrape');
    console.log("scrape response", resp);
    const data = await resp.json();
    setScrapedBlogPosts(data);
  };

  const newBlogPostTitle = useInput("");
  const newBlogPostURL = useInput("");
  const [savedBlogPosts, setSavedBlogPosts] = useState([]);
  const addAndSaveButtonOnPress = async (_: PressEvent) => {
    const payload: Array<BlogPost> = Object.assign([], scrapedBlogPosts);
    payload.push({
      title: newBlogPostTitle.value,
      url: newBlogPostURL.value,
    });

    const resp = await fetch('http://localhost:8000/api/persist-first', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });
    console.log("persist first response", resp);
    const data = await resp.json();
    setSavedBlogPosts(data);
  };

  return (
    <main>
      <Button onPress={scrapeButtonOnPress}>Scrape Drew DeVault's blog</Button>

      {
        blogPostsTableFC('Scraped blog posts', scrapedBlogPosts)
      }

      <Spacer y={2} />

      <Input label="Title" placeholder="Title" {...newBlogPostTitle.bindings}></Input>
      <Spacer y={1} />
      <Input label="URL" placeholder="URL" {...newBlogPostURL.bindings}></Input>
      <Spacer y={1} />
      <Button onPress={addAndSaveButtonOnPress}>Add new blog post and save all</Button>

      {
        blogPostsTableFC('Saved blog posts', savedBlogPosts)
      }
    </main>
  )
}
