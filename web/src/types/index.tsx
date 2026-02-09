export type Repository = {
  id: string;
  name: string;
  owner: string;
  full_name?: string;
  description?: string;
  language?: string;
  stars?: number;
  forks?: number;
  openIssues?: number;
  lastUpdated?: string;
  updated_at?: string;
  url?: string;
  html_url?: string;
  private?: boolean;
  is_private?: boolean;
};

export type LanguageColor = {
  [language: string]: string; // Maps a language to its hex color
};

export type User = {
  id: string;
  github_id: string;
  username: string;
  email?: string;
  avatar_url?: string;
  access_token?: string;
  created_at: string;
  updated_at: string;
};
