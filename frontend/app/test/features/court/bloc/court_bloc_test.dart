// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:bloc_test/bloc_test.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:qlove/core/models/court_case_model.dart';
import 'package:qlove/features/court/bloc/court_bloc.dart';
import 'package:qlove/features/court/bloc/court_event.dart';
import 'package:qlove/features/court/bloc/court_state.dart';
import 'package:qlove/features/court/data/court_repository.dart';

class MockCourtRepository extends Mock implements CourtRepository {}

void main() {
  late CourtBloc courtBloc;
  late MockCourtRepository mockCourtRepository;

  setUp(() {
    mockCourtRepository = MockCourtRepository();
    courtBloc = CourtBloc(courtRepository: mockCourtRepository);
  });

  tearDown(() {
    courtBloc.close();
  });

  group('CourtBloc', () {
    final mockCases = [
      CourtCaseModel(id: '1', defendantNameMasked: 'A ***', reason: 'Ghosting', status: 'voting', voteCount: 10, createdAt: DateTime.now()),
      CourtCaseModel(id: '2', defendantNameMasked: 'B ***', reason: 'Trap', status: 'voting', voteCount: 5, createdAt: DateTime.now()),
    ];

    test('initial state is CourtInitial', () {
      expect(courtBloc.state, isA<CourtInitial>());
    });

    blocTest<CourtBloc, CourtState>(
      'emits [CourtLoading, CourtLoaded] when FetchCasesRequested succeeds',
      build: () {
        when(() => mockCourtRepository.getCases())
            .thenAnswer((_) async => mockCases);
        return courtBloc;
      },
      act: (bloc) => bloc.add(const FetchCasesRequested()),
      expect: () => [
        isA<CourtLoading>(),
        isA<CourtLoaded>().having((state) => state.cases.length, 'cases length', 2),
      ],
    );

    blocTest<CourtBloc, CourtState>(
      'emits [CourtLoaded] optimally when VoteActionRequested succeeds',
      build: () {
        when(() => mockCourtRepository.vote('1', 'guilty'))
            .thenAnswer((_) async => Future.value());
        return courtBloc;
      },
      seed: () => CourtLoaded(cases: mockCases, hasReachedMax: false),
      act: (bloc) => bloc.add(const VoteActionRequested(caseId: '1', voteType: 'guilty')),
      expect: () => [
        isA<CourtLoaded>().having((state) => state.cases.length, 'cases length after vote', 1),
      ],
    );
  });
}
