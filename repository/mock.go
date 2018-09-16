package repository

// Mock is mocking the Repository interface for testing
type Mock struct {
	StoredGroup  Group
	StoredConfig Config
}

// StoreGroup stores Group object to the mock
func (r *Mock) StoreGroup(g Group) error {
	r.StoredGroup = g

	return nil
}

// RetrieveGroup retrieves Group object from the mock
func (r *Mock) RetrieveGroup(id string) (*Group, error) {
	if id != r.StoredGroup.ID {
		return nil, ErrGroupNotFound
	}

	return &r.StoredGroup, nil
}

// StoreConfig stores Config object to the mock
func (r *Mock) StoreConfig(c Config) error {
	grp, err := r.RetrieveGroup(c.Group)
	if err != nil {
		return err
	}

	r.StoredConfig = c

	grp.Configs = append(grp.Configs, c.ID)
	err = r.StoreGroup(*grp)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveConfig retrieves Config object from the mock
func (r *Mock) RetrieveConfig(groupID string, id string) (*Config, error) {
	if groupID != r.StoredGroup.ID {
		return nil, ErrGroupNotFound
	}

	if id != r.StoredConfig.ID {
		return nil, ErrConfigNotFound
	}

	return &r.StoredConfig, nil
}
