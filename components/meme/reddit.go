/*
 * JOJO Discord Bot - An advanced multi-purpose discord bot
 * Copyright (C) 2022 Lazy Bytez (Elias Knodel, Pascal Zarrad)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package meme

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/build"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"math/rand"
)

// RedditMeme is a struct that represents a meme from reddit.
type RedditMeme struct {
	// The title of the meme
	Title string `json:"title"`
	// The URL of the meme
	URL string `json:"url"`
	// The subreddit of the meme
	Subreddit string `json:"subreddit"`
	// The Image URL of the meme
	ImageURL string `json:"image_url"`
}

const (
	// ErrInvalidMemeSubreddit is the error message that is returned when the user
	// specifies an invalid subreddit.
	ErrInvalidMemeSubreddit = "specified subreddit is not in the whitelist"
	// ErrRedditAPIError is the error message that is returned when an error occurs
	// while fetching the meme from reddit.
	ErrRedditAPIError = "an error occurred while fetching the meme from reddit"
	// ErrorMemeNotFound is the error message that is returned when no meme could be
	// found in the specified subreddit.
	ErrorMemeNotFound = "no meme found in the subreddit %s"
)

// RedditFindMemeTries is the amount of tries to find a meme in the specified subreddit
// before giving up and returning an error.
const RedditFindMemeTries = 3

// randomMemeSubredditCollection is a collection of meme subreddits.
// The user can choose one of these subreddits to get a random meme from.
var memeSubreddits = []*discordgo.ApplicationCommandOptionChoice{
	{
		Name:  "r/ProgrammerHumor",
		Value: "ProgrammerHumor",
	},
	{
		Name:  "r/mememes",
		Value: "mememes",
	},
	{
		Name:  "r/dankmemes",
		Value: "dankmemes",
	},
	{
		Name:  "r/wholesomememes",
		Value: "wholesomememes",
	},
	{
		Name:  "r/meirl",
		Value: "meirl",
	},
	{
		Name:  "r/ich_iel",
		Value: "ich_iel",
	},
	{
		Name:  "r/animemes",
		Value: "animemes",
	},
}

// randomMemeSubredditCollection is a collection of meme subreddits
// that are used to randomly select a subreddit.
var randomMemeSubredditCollection []string

// init initializes the list of random meme subreddits
func init() {
	initializeRandomMemeSubreddits()
}

// initializeRandomMemeSubreddits initializes the list of random meme subreddits
// by iterating over the memeSubreddits slice.
func initializeRandomMemeSubreddits() {
	for _, choice := range memeSubreddits {
		randomMemeSubredditCollection = append(randomMemeSubredditCollection, choice.Value.(string))
	}
}

// findRandomMeme finds a random meme from the given subreddit.
// It will try to find a meme up to RedditFindMemeTries times.
func findRandomMeme(subreddit string, tries int) (RedditMeme, error) {
	meme, err := getRandomMeme(subreddit)
	if err == nil {
		return meme, nil
	}

	if tries < RedditFindMemeTries {
		return findRandomMeme(subreddit, tries+1)
	}

	return meme, err
}

// randomMemeSubreddit returns a random meme either from the specified subreddit
// or from a random subreddit.
func getRandomMeme(subreddit string) (RedditMeme, error) {
	if subreddit == "" {
		subreddit = randomMemeSubreddit()
	}

	if !isMemeSubreddit(subreddit) {
		return RedditMeme{}, fmt.Errorf(ErrInvalidMemeSubreddit)
	}

	botUa := fmt.Sprintf("jojo-discord-bot/%s", build.ComputeVersionString())
	redditApi, err := reddit.NewClient(reddit.Credentials{
		ID:       "YOUR_CLIENT_ID",
		Secret:   "YOUR_SECRET",
		Username: "YOUR_REDDIT_USERNAME",
		Password: "YOUR_REDDIT_PASSWORD",
	}, reddit.WithUserAgent(botUa))

	if err != nil {
		return RedditMeme{}, fmt.Errorf(ErrRedditAPIError)
	}

	posts, _, err := redditApi.Subreddit.TopPosts(context.Background(), subreddit, &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 100,
		},
	})

	if len(posts) == 0 {
		return RedditMeme{}, fmt.Errorf(ErrorMemeNotFound, subreddit)
	}

	// Choose a random post from the top posts.
	post := posts[rand.Intn(len(posts))]
	if post.IsSelfPost {
		return RedditMeme{}, fmt.Errorf(ErrorMemeNotFound, subreddit)
	}

	return RedditMeme{
		Title:     post.Title,
		URL:       post.URL,
		Subreddit: post.SubredditNamePrefixed,
		ImageURL:  post.URL,
	}, nil
}

// isMemeSubreddit checks if the specified subreddit is in the list of meme subreddits.
func isMemeSubreddit(meme string) bool {
	for _, choice := range memeSubreddits {
		if choice.Value.(string) == meme {
			return true
		}
	}

	return false
}

// randomMemeSubreddit returns a random meme subreddit from the randomMemeSubredditCollection.
func randomMemeSubreddit() string {
	return randomMemeSubredditCollection[rand.Intn(len(randomMemeSubredditCollection))]
}
