# Feature Flag
This helper will allow you to use Flagr as Feature Flag service. Read more about Feature Flag in [here](https://martinfowler.com/articles/feature-toggles.html) and Flagr in [here](https://checkr.github.io/flagr/#/).

## Usage
Create `FeatureFlagConfig` first then call `NewFeatureFlag` by injecting config and context into it. Note : Make sure you append `/api/v1` in the `BaseURL`.
```go
    config := &FeatureFlagConfig{
		BaseURL:  "http://some-url/api/v1",
		Username: "someuser",
		Password: "somepassword",
	}
	ctx := context.Background()
	featureFlag := NewFeatureFlag(config, ctx)
```

### Evaluate
Evaluate by variant and get the distribution result.
```go
    ok, err := featureFlag.Eval('some-flag', 'some-variant')
```

### Evaluate then execute function
Evaluate by variant and execute function if true.
```go
    ok, err := featureFlag.EvalThenExecute('some-flag', 'some-variant', func() error {
        //Do something
        return err
    })
```