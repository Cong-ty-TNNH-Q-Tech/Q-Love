import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:rive/rive.dart';
import 'island_state.dart';

class IslandCubit extends Cubit<IslandState> {
  IslandCubit() : super(IslandInitial());

  StateMachineController? controller;

  void onRiveInit(Artboard artboard) {
    controller = StateMachineController.fromArtboard(artboard, 'IslandStateMachine');
    if (controller != null) {
      artboard.addController(controller!);
    }
    emit(IslandLoaded());
  }

  @override
  Future<void> close() {
    controller?.dispose();
    return super.close();
  }
}
