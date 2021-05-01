package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTagService_GetModel(t *testing.T) {
	setupTest(t)

	node := Node{
		Name:     "test node",
		IsMaster: true,
		Tags: []Tag{
			{Name: "tag 1", Color: "red"},
		},
	}
	err := node.Add()
	require.Nil(t, err)

	tag, err := MustGetRootService().GetTag(nil, nil)
	require.Nil(t, err)
	require.False(t, tag.Id.IsZero())
	require.Equal(t, "tag 1", tag.Name)
	require.Equal(t, "red", tag.Color)
	require.Equal(t, interfaces.ModelColNameNode, tag.Col)
}

func TestTagService_GetModelById(t *testing.T) {
	setupTest(t)

	node := Node{
		Name:     "test node",
		IsMaster: true,
		Tags: []Tag{
			{Name: "tag 1", Color: "red"},
		},
	}
	err := node.Add()
	require.Nil(t, err)

	tag, err := MustGetRootService().GetTag(nil, nil)
	require.Nil(t, err)
	require.False(t, tag.Id.IsZero())
	require.Equal(t, "tag 1", tag.Name)
	require.Equal(t, "red", tag.Color)
	require.Equal(t, interfaces.ModelColNameNode, tag.Col)
}
