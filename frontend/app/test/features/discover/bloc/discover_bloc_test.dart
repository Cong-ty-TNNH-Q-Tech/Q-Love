// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:bloc_test/bloc_test.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:qlove/core/models/user_model.dart';
import 'package:qlove/features/discover/bloc/discover_bloc.dart';
import 'package:qlove/features/discover/bloc/discover_event.dart';
import 'package:qlove/features/discover/bloc/discover_state.dart';
import 'package:qlove/features/discover/data/discover_repository.dart';

class MockDiscoverRepository extends Mock implements DiscoverRepository {}

void main() {
  late DiscoverBloc discoverBloc;
  late MockDiscoverRepository mockDiscoverRepository;

  setUp(() {
    mockDiscoverRepository = MockDiscoverRepository();
    discoverBloc = DiscoverBloc(discoverRepository: mockDiscoverRepository);
  });

  tearDown(() {
    discoverBloc.close();
  });

  group('DiscoverBloc', () {
    final mockProfiles = [
      const UserModel(id: '1', name: 'Alice'),
      const UserModel(id: '2', name: 'Bob'),
    ];

    test('initial state is DiscoverInitial', () {
      expect(discoverBloc.state, isA<DiscoverInitial>());
    });

    blocTest<DiscoverBloc, DiscoverState>(
      'emits [DiscoverLoading, DiscoverLoaded] when FetchFeedRequested succeeds',
      build: () {
        when(() => mockDiscoverRepository.getFeed(filter: 'default'))
            .thenAnswer((_) async => mockProfiles);
        return discoverBloc;
      },
      act: (bloc) => bloc.add(const FetchFeedRequested()),
      expect: () => [
        isA<DiscoverLoading>(),
        isA<DiscoverLoaded>().having((state) => state.profiles.length, 'profiles length', 2),
      ],
    );

    blocTest<DiscoverBloc, DiscoverState>(
      'emits [DiscoverLoaded] optimally and then [DiscoverMatch] when SwipeActionRequested is like and matches',
      build: () {
        when(() => mockDiscoverRepository.getFeed(filter: 'default'))
            .thenAnswer((_) async => mockProfiles);
        when(() => mockDiscoverRepository.swipe('1', 'like'))
            .thenAnswer((_) async => true);
        return discoverBloc;
      },
      seed: () => DiscoverLoaded(profiles: mockProfiles, hasReachedMax: false),
      act: (bloc) => bloc.add(const SwipeActionRequested(targetId: '1', action: 'like')),
      expect: () => [
        isA<DiscoverLoaded>().having((state) => state.profiles.length, 'profiles length after swipe', 1),
        isA<DiscoverMatch>().having((state) => state.matchedUser.id, 'matched user id', '1'),
      ],
    );
  });
}
