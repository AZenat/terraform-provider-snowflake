package sdk

import "testing"

func TestApplicationPackages_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *CreateApplicationPackageOptions {
		return &CreateApplicationPackageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateApplicationPackageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.DataRetentionTimeInDays = Int(1)
		opts.MaxDataExtensionTimeInDays = Int(1)
		opts.DefaultDdlCollation = String("en_US")
		opts.Comment = String("comment")
		opts.Distribution = DistributionPointer(DistributionInternal)
		t1 := randomSchemaObjectIdentifier()
		opts.Tag = []TagAssociation{
			{
				Name:  t1,
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE APPLICATION PACKAGE IF NOT EXISTS %s DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 DEFAULT_DDL_COLLATION = 'en_US' COMMENT = 'comment' DISTRIBUTION = INTERNAL TAG (%s = 'v1')", id.FullyQualifiedName(), t1.FullyQualifiedName())
	})
}

func TestApplicationPackages_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *AlterApplicationPackageOptions {
		return &AlterApplicationPackageOptions{
			IfExists: Bool(true),
			name:     id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterApplicationPackageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApplicationPackageOptions", "Set", "Unset", "ModifyReleaseDirective", "SetDefaultReleaseDirective", "SetReleaseDirective", "UnsetReleaseDirective", "AddVersion", "DropVersion", "AddPatchForVersion", "SetTags", "UnsetTags"))
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetDefaultReleaseDirective = &SetDefaultReleaseDirective{
			Version: "v1",
			Patch:   1,
		}
		opts.UnsetReleaseDirective = &UnsetReleaseDirective{
			ReleaseDirective: "DEFAULT",
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApplicationPackageOptions", "Set", "Unset", "ModifyReleaseDirective", "SetDefaultReleaseDirective", "SetReleaseDirective", "UnsetReleaseDirective", "AddVersion", "DropVersion", "AddPatchForVersion", "SetTags", "UnsetTags"))
	})

	t.Run("validation: set options at least one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApplicationPackageUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApplicationPackageOptions.Unset", "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "DefaultDdlCollation", "Comment", "Distribution"))
	})

	t.Run("alter: set options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApplicationPackageSet{
			DataRetentionTimeInDays:    Int(1),
			MaxDataExtensionTimeInDays: Int(1),
			DefaultDdlCollation:        String("en_US"),
			Comment:                    String("comment"),
			Distribution:               DistributionPointer(DistributionInternal),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s SET DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 DEFAULT_DDL_COLLATION = 'en_US' COMMENT = 'comment' DISTRIBUTION = INTERNAL`, id.FullyQualifiedName())
	})

	t.Run("alter: unset options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApplicationPackageUnset{
			DataRetentionTimeInDays:    Bool(true),
			MaxDataExtensionTimeInDays: Bool(true),
			DefaultDdlCollation:        Bool(true),
			Comment:                    Bool(true),
			Distribution:               Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, DEFAULT_DDL_COLLATION, COMMENT, DISTRIBUTION`, id.FullyQualifiedName())
	})

	t.Run("alter: set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s SET TAG "tag1" = 'value1'`, id.FullyQualifiedName())
	})

	t.Run("alter: unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})

	t.Run("alter: modify release directive", func(t *testing.T) {
		opts := defaultOpts()
		opts.ModifyReleaseDirective = &ModifyReleaseDirective{
			ReleaseDirective: "DEFAULT",
			Version:          "V1",
			Patch:            1,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s MODIFY RELEASE DIRECTIVE DEFAULT VERSION = V1 PATCH = 1`, id.FullyQualifiedName())
	})

	t.Run("alter: set default release directive", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetDefaultReleaseDirective = &SetDefaultReleaseDirective{
			Version: "V1",
			Patch:   1,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s SET DEFAULT RELEASE DIRECTIVE VERSION = V1 PATCH = 1`, id.FullyQualifiedName())
	})

	t.Run("alter: set release directive", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetReleaseDirective = &SetReleaseDirective{
			ReleaseDirective: "DEFAULT",
			Accounts: []string{
				"org1.acc1",
				"org2.acc2",
			},
			Version: "V1",
			Patch:   1,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s SET RELEASE DIRECTIVE DEFAULT ACCOUNTS = (org1.acc1, org2.acc2) VERSION = V1 PATCH = 1`, id.FullyQualifiedName())
	})

	t.Run("alter: set release directive with no accounts", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetReleaseDirective = &SetReleaseDirective{
			ReleaseDirective: "DEFAULT",
			Version:          "V1",
			Patch:            1,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s SET RELEASE DIRECTIVE DEFAULT ACCOUNTS = () VERSION = V1 PATCH = 1`, id.FullyQualifiedName())
	})

	t.Run("alter: unset release directive", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetReleaseDirective = &UnsetReleaseDirective{
			ReleaseDirective: "DEFAULT",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s UNSET RELEASE DIRECTIVE DEFAULT`, id.FullyQualifiedName())
	})

	t.Run("alter: add version", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddVersion = &AddVersion{
			VersionIdentifier: String("v1_1"),
			Using:             "@hello_snowflake_code.core.hello_snowflake_stage",
			Label:             String("test"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s ADD VERSION v1_1 USING '@hello_snowflake_code.core.hello_snowflake_stage' LABEL = 'test'`, id.FullyQualifiedName())
	})

	t.Run("alter: drop version", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropVersion = &DropVersion{
			VersionIdentifier: "v1_1",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s DROP VERSION v1_1`, id.FullyQualifiedName())
	})

	t.Run("alter: add patch for version", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddPatchForVersion = &AddPatchForVersion{
			VersionIdentifier: String("v1_1"),
			Using:             "@hello_snowflake_code.core.hello_snowflake_stage",
			Label:             String("test"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION PACKAGE IF EXISTS %s ADD PATCH FOR VERSION v1_1 USING '@hello_snowflake_code.core.hello_snowflake_stage' LABEL = 'test'`, id.FullyQualifiedName())
	})
}

func TestApplicationPackages_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *DropApplicationPackageOptions {
		return &DropApplicationPackageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropApplicationPackageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP APPLICATION PACKAGE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestApplicationPackages_Show(t *testing.T) {
	defaultOpts := func() *ShowApplicationPackageOptions {
		return &ShowApplicationPackageOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowApplicationPackageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW APPLICATION PACKAGES`)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		opts.StartsWith = String("A")
		opts.Limit = &LimitFrom{
			Rows: Int(1),
			From: String("B"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW APPLICATION PACKAGES LIKE 'pattern' STARTS WITH 'A' LIMIT 1 FROM 'B'`)
	})
}
